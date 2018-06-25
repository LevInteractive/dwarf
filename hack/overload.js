const request = require("request-promise-native");

function count(needle, haystack) {
  let count = 0;
  for (let i = 0; i < haystack.length; i++) {
    if (needle === haystack[i]) {
      count++;
    }
  }
  return count;
}

async function main() {
  const total = 1000;
  const urls = [];
  for (let i = 0; i < total; i++) {
    urls.push(`https://google.com/${i}`);
  }
  const reqs = [];

  // TEST BULK
  const opts = {
    method: "POST",
    uri: "http://localhost:3001/create",
    body: {
      longUrl: urls,
      apiKey: "123"
    },
    json: true
  };
  const results = await request.post(opts);
  console.log(results);

  if (results.length !== total) {
    throw new Error("Did not equal length 1000. got " + results.length);
  }
  /////////////////////////////////////////////

  // TEST SINGLE
  //
  // for (let i = 0; i < total; i++) {
  //   const opts = {
  //     method: "POST",
  //     uri: "http://localhost:3001/create",
  //     body: {
  //       longUrl: urls[i],
  //       apiKey: "123"
  //     },
  //     json: true
  //   };
  //   reqs.push(request.post(opts));
  // }

  // const results = await Promise.all(reqs);
  // const shortUrls = results.map(r => r.shortUrl);
  //
  // for (let i = 0; i < shortUrls.length; i++) {
  //   if (count(shortUrls[i], shortUrls) !== 1) {
  //     throw new Error("Not equal to one: ", shortUrls);
  //   }
  // }
}

main();
