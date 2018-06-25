const http = require("http");
const bodyParser = require("body-parser");
const express = require("express");
const { connect, UrlShort } = require("./utils/mongo");
const config = require("./utils/config");
const { log, error } = require("./utils/logging");
const app = express();
const cors = require("cors");

// Because it's running behind nginx in prod.
// https://github.com/expressjs/morgan/issues/114
app.enable("trust proxy");
app.set("trust proxy", () => true);

// Setup bodyParser with catch JSON errors
app.use(bodyParser.json());
app.use(async (err, req, res, next) => {
  if (err instanceof SyntaxError && err.status === 400) {
    error("Bad JSON", req.body);
    return res.status(404).json({ error: true, message: "Bad JSON" });
  }
  if (err) {
    return res.status(404).json({ error: true, message: err.message });
  }
});

// Cors configuration
const whitelist = ["http://localhost:3000", "http://(.*).localhost:3000"];
const corsOptions = { origin: true };
if (config.whitelist && config.whitelist !== "*") {
  const whitelist = config.whitelist.split(",");
  whitelist.forEach(function(val, idx) {
    if (whitelist[idx].match(/^\/(.*)\/$/)) {
      whitelist[idx] = new RegExp(
        whitelist[idx]
          .substr(1, whitelist[idx].length - 2)
          .replace("\\\\", "\\")
      );
    }
  });
  corsOptions.origin = whitelist;
}
app.use(cors(corsOptions));
app.options("*", cors(corsOptions));

app.post("/create", async (req, res) => {
  try {
    log(`[REQUEST FROM] ${req.hostname}`);
    if (
      !req.body.apiKey ||
      (req.body.apiKey && req.body.apiKey !== config.apiKey)
    ) {
      return res
        .status(404)
        .json({ error: true, message: "You need to send a valid apiKey" });
    }
    if (!req.body.longUrl) {
      return res.status(404).json({
        error: true,
        message: "You need to send a longUrl to be shorten"
      });
    }
    const longUrl = req.body.longUrl;
    const result = await UrlShort.shorten(longUrl, req.body.code);

    return res.status(200).json(result);
  } catch (e) {
    return res.status(404).json({ error: true, message: e.message });
  }
});

app.get("/:code", async (req, res) => {
  const code = req.params.code;
  const urlShort = await UrlShort.findByCode(code);
  if (!urlShort) {
    log(`[NOT FOUND] ${config.baseUrl}/${code}`);
    return res.status(404).json({ error: "Not found" });
  } else {
    log(`[REDIRECTING] ${config.baseUrl}/${code} => ${urlShort.longUrl}`);
    return res.redirect(301, urlShort.longUrl);
  }
});

app.get("/", (req, res) => {
  return res.send(
    `${
      config.baseUrl
    } <a href="https://github.com/LevInteractive/dwarf">DWARF shortener</a>`
  );
});

app.set("port", config.port);

(async function main() {
  await connect();
  log(`DWARF Url Shortener running on host ${config.baseUrl}`);
  const server = http.createServer(app);
  server.listen(config.port);
})();
