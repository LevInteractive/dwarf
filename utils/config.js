const config = require("dotenv").config();

// exports.mongoConnectionString =
//   config.parsed.MONGO_CONNECTION_STRING || "mongodb://localhost/dwarf";
exports.mongoConnectionString =
  config.parsed.MONGO_CONNECTION_STRING || "mongodb://127.0.0.1:27017";
exports.mongoDatabase = config.parsed.MONGO_DATABASE || "dwarf";
exports.port = config.parsed.PORT || 3001;
exports.baseUrl = config.parsed.BASE_URL || "http://localhost:3001";
exports.minChars = config.parsed.MIN_CHARS || 0;
exports.apiKey = config.parsed.API_KEY || "CHANGE_API_KEY";
