declare module "streamsaver" {
  export function createWriteStream(filename: string, options?: { size?: number }): WritableStream;
  export default { createWriteStream };
}
