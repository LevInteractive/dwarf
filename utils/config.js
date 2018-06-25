const config = require("dotenv").config();

const Config = {
  redisHost: "127.0.0.1",
  redisPort: "6379",
  mongoConnectionString: "mongodb://127.0.0.1:27017",
  mongoDatabase: "dwarf",
  port: 3001,
  baseUrl: "http://localhost:3001",
  apiKey: "CHANGE_API_KEY",
  whitelist: "*"
};

const params = {
  redisHost: "REDIS_HOST",
  redisPort: "REDIS_PORT",
  mongoConnectionString: "MONGO_CONNECTION_STRING",
  mongoDatabase: "MONGO_DATABASE",
  port: "PORT",
  baseUrl: "BASE_URL",
  apiKey: "API_KEY",
  whitelist: "WHITELIST"
};

Object.keys(params).forEach(function(idx) {
  const key = params[idx];
  // First get .env configuration if available
  if (!config.error) {
    if (config.parsed[key]) {
      Config[idx] = config.parsed[key];
    }
  }
  // Overrides with process.env vars if sent
  if (process.env[key]) {
    Config[idx] = process.env[key];
  }
});

module.exports = Config;
