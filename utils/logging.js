/**
 * Log wrapper
 *
 * @module utils/logging
 */
const winston = require("winston");

/**
 * Console log - info level
 *
 * @return {function}
 */
exports.info = exports.log = function() {
  if (process.env.NODE_ENV !== "test") {
    winston.info.apply(winston, arguments);
  }
};

/**
 * Console log - error level
 *
 * @return {function}
 */
exports.error = winston.error;

/**
 * Console log - warn level
 *
 * @return {function}
 */
exports.warn = function() {
  if (process.env.NODE_ENV !== "test") {
    winston.warn.apply(winston, arguments);
  }
};

/**
 * Time Start
 *
 * @param  {string} str Timer name
 * @return {void}
 */
exports.timeStart = str => {
  if (process.env.NODE_ENV !== "test") {
    console.time(str);
  }
};

/**
 * Time End
 *
 * @param  {string} str Timer name
 * @return {void}
 */
exports.timeEnd = str => {
  if (process.env.NODE_ENV !== "test") {
    console.timeEnd(str);
  }
};
