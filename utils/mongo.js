/**
 * MongoDB connection
 *
 * @module utils/mongo
 *
 * @requires NPM:mongodb
 * @requires NPM:pluralize
 * @requires NPM:config
 * @requires ./logging
 */

const { MongoClient } = require("mongodb");
const { plural } = require("pluralize");
const config = require("./config");
const { log, error } = require("./logging");
const redis = require("redis");

const redisPrefix = "dwarf:";

const redisClient = redis.createClient({
  host: config.redisHost,
  port: config.redisPort
});

/**
 * Regex to check MongoID
 *
 * @type {RegExp}
 */
const checkForHexRegExp = new RegExp("^[0-9a-fA-F]{24}$");

/**
 * Database persistent
 *
 * @type {Object}
 */
let _db;

/**
 * Connect to the database
 *
 * @return {Promise}
 */
exports.connect = function() {
  return new Promise((resolve, reject) => {
    if (_db) {
      return resolve();
    }

    const opts = {
      keepAlive: 1,
      connectTimeoutMS: 30000
    };

    MongoClient.connect(
      config.mongoConnectionString,
      function(err, conn) {
        if (err) {
          reject(err);
          return error(
            `Failed to connect to MongoDB: ${config.mongoConnectionString}/${
              config.mongoDatabase
            }`
          );
        }
        _db = conn.db(config.mongoDatabase);
        log(
          `Successfully connected to MongoDB: ${config.mongoConnectionString}/${
            config.mongoDatabase
          }`
        );
        resolve(_db);
      }
    );
  });
};

/**
 * Prime the counter based on number of existing URL's in the database. This
 * is a temporary hack until we move to 100% redis.
 */
exports.init = async function() {
  const model = modelQuery("UrlShort");
  const currentCount = await model.count();
  redisClient.set(redisPrefix + "counter", currentCount);
};

/**
 * Connect to mongo.
 *
 * @type {Object}
 */
const db = (exports.db = () => _db);

/**
 * Return a query object for the database.
 *
 * @param  {string} singularName Singular name of the collection/provider
 * @return {Object}
 */
function modelQuery(singularName) {
  if (typeof singularName !== "string") {
    const err = new Error(
      "You must pass a string to modelQuery. Passed: " + singularName
    );
    console.error(err.stack);
    throw err;
  }
  const collectionName = plural(singularName.toLowerCase());

  return db().collection(collectionName);
}
exports.modelQuery = modelQuery;

const getRandomInt = function(min, max) {
  return Math.floor(Math.random() * (max - min)) + min;
};

async function getCounter() {
  return new Promise((resolve, reject) => {
    redisClient.incr(redisPrefix + "counter", (err, count) => {
      if (err) {
        reject(err);
      } else {
        resolve(getRandomInt(9999, 999999) + count.toString());
      }
    });
  });
}

exports.UrlShort = {
  shorten: async function(longUrl, code) {
    const model = modelQuery("UrlShort");
    let count;

    if (Array.isArray(longUrl)) {
      const urls = [];

      for (let i = 0, len = longUrl.length; i < len; i++) {
        const lUrl = longUrl[i];
        if (typeof lUrl !== "string") {
          error(`[ERROR: LONG_URL is not a string] ${lUrl}`);
          return {
            longUrl: lUrl,
            error: true,
            message: "longUrl is not a string"
          };
        }

        if (!isValidUrl(lUrl)) {
          error(`[ERROR: LONG_URL is an invalid URL] ${lUrl}`);
          return {
            longUrl: lUrl,
            error: true,
            message:
              "Invalid URL format. Input URL must comply to the following: http(s)://(www.)domain.ext(/)(path)"
          };
        }

        // Check if longUrl already exists, so just the shortUrl
        const existingUrl = await model.findOne({ longUrl: lUrl });
        if (existingUrl) {
          log(
            `[ABORTING: LONG_URL already exists] ${lUrl} => ${config.baseUrl}/${
              existingUrl.code
            }`
          );

          urls.push({
            longUrl: lUrl,
            shortUrl: `${config.baseUrl}/${existingUrl.code}`
          });
        }
        // Doesn't exist, let's create
        count = await getCounter();
        const code = encode(count);
        const fields = {
          _id: count,
          longUrl: lUrl,
          code,
          created: new Date()
        };
        await model.insertOne(fields);
        log(`[CREATED] ${lUrl} => ${config.baseUrl}/${code}`);
        urls.push({ longUrl: lUrl, shortUrl: `${config.baseUrl}/${code}` });
      }

      return urls;
    }

    if (typeof longUrl !== "string") {
      throw new Error("longUrl is not a string");
    }

    if (!isValidUrl(longUrl)) {
      throw new Error(
        "Invalid URL format. Input URL must comply to the following: http(s)://(www.)domain.ext(/)(path)"
      );
    }

    // Check if code is sent, no need to generate
    if (code) {
      // Check if code already exists, so just return shortUrl
      const existingUrl = await model.findOne({ code });
      if (existingUrl) {
        log(
          `[ABORTING: CODE already exists] ${longUrl} => ${config.baseUrl}/${
            existingUrl.code
          }`
        );
        return { longUrl, shortUrl: `${config.baseUrl}/${existingUrl.code}` };
      }
      count = await getCounter();
    } else {
      // Check if longUrl already exists, so just return shortUrl
      const existingUrl = await model.findOne({ longUrl });
      if (existingUrl) {
        log(
          `[ABORTING: LONG_URL already exists] ${longUrl} => ${
            config.baseUrl
          }/${existingUrl.code}`
        );

        return { longUrl, shortUrl: `${config.baseUrl}/${existingUrl.code}` };
      }
      // Nope, just generate a incremented one
      count = await getCounter();
      code = encode(count);
    }

    const fields = {
      _id: count,
      longUrl,
      code,
      created: new Date()
    };
    await model.insertOne(fields);
    log(`[CREATED] ${longUrl} => ${config.baseUrl}/${code}`);

    return { longUrl, shortUrl: `${config.baseUrl}/${code}` };
  },

  findByCode: async function(code) {
    const model = modelQuery("UrlShort");
    return await model.findOne({ code });
  }
};

const alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ";
const base = alphabet.length; // base is the length of the alphabet (58 in this case)

// Utility function to convert base 10 integer to base 58 string
function encode(num) {
  let encoded = "";
  while (num) {
    const remainder = num % base;
    num = Math.floor(num / base);
    encoded = alphabet[remainder].toString() + encoded;
  }
  return encoded;
}

// Utility function to convert a base 58 string to base 10 integer
function decode(str) {
  let decoded = 0;
  while (str) {
    const index = alphabet.indexOf(str[0]);
    const power = str.length - 1;
    decoded += index * Math.pow(base, power);
    str = str.substring(1);
  }
  return decoded;
}

function isValidUrl(url) {
  // Must comply to this format () means optional:
  // http(s)://(www.)domain.ext(/)(whatever follows)
  const regEx = /^https?:\/\/(\S+\.)?(\S+\.)(\S+)\S*/;
  return regEx.test(url);
}
