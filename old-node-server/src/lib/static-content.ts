import pathlib from "path";
import * as ExpressCore from "express-serve-static-core";

export const registerStaticRequestHandlers = (expressApp: ExpressCore.Express) => {

  expressApp.get("/static/stream-saver/mitm_2_0_0.html", (req, res) => {
    let description = `STATIC ${req.method} ${req.url}`;
    logger.log(description);
    let localFilePath = pathlib.resolve('./src/static/stream-saver/mitm_2_0_0.html')
    return res.sendFile(localFilePath);
  });

  expressApp.get("/static/stream-saver/sw_2_0_0.js", (req, res) => {
    let description = `STATIC ${req.method} ${req.url}`;
    logger.log(description);
    let localFilePath = pathlib.resolve('./src/static/stream-saver/sw_2_0_0.js')
    return res.sendFile(localFilePath);
  });

}