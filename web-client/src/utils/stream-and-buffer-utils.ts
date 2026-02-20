export const convertStreamToBuffer = async (readableStream: ReadableStream): Promise<ArrayBuffer> => {
  const arrayBufferList: ArrayBuffer[] = [];

  const reader = readableStream.getReader();
  while (true) {
    const { done, value: chunk } = await reader.read();
    await new Promise((accept) => setTimeout(accept, 40));
    if (done) break;

    let arrayBuffer: ArrayBuffer;
    if (chunk instanceof Uint8Array) {
      arrayBuffer = chunk.buffer.slice(chunk.byteOffset, chunk.byteOffset + chunk.byteLength) as ArrayBuffer;
    } else {
      arrayBuffer = chunk as ArrayBuffer;
    }
    arrayBufferList.push(arrayBuffer);
  }

  let totalLength = 0;
  arrayBufferList.forEach((arrayBuffer) => {
    const view = new Uint8Array(arrayBuffer);
    totalLength += view.length;
  });

  const resultBuffer = new ArrayBuffer(totalLength);
  const resultView = new Uint8Array(resultBuffer);

  let startIndex = 0;
  for (const arrayBuffer of arrayBufferList) {
    const view = new Uint8Array(arrayBuffer);
    resultView.set(view, startIndex);
    startIndex += view.length;
  }

  return resultBuffer;
};

export class MeteredByteStreamReader {
  private readableStream: ReadableStream;
  private reader: ReadableStreamDefaultReader;

  private remainingChunkOffset = 0;
  private remainingChunk: Uint8Array | null = null;
  id: string;

  constructor(readableStream: ReadableStream, id: string) {
    this.readableStream = readableStream;
    this.reader = this.readableStream.getReader();
    this.id = id;
  }

  async readBytes(byteCount: number): Promise<{ value: Uint8Array; done: boolean }> {
    let returnBytes = new Uint8Array(byteCount);
    let startingOffset = 0;
    let isDone = false;

    if (this.remainingChunk && this.remainingChunk.length === this.remainingChunkOffset) {
      this.remainingChunk = null;
      this.remainingChunkOffset = 0;
    }

    if (this.remainingChunk && this.remainingChunk.length > this.remainingChunkOffset) {
      const byteCountToCopy = Math.min(byteCount, this.remainingChunk.length - this.remainingChunkOffset);

      const sourceArray = this.remainingChunk.slice(this.remainingChunkOffset, this.remainingChunkOffset + byteCountToCopy);

      returnBytes.set(sourceArray, startingOffset);
      startingOffset += byteCountToCopy;
      this.remainingChunkOffset += byteCountToCopy;
    }

    while (true) {
      if (startingOffset === byteCount) break;

      const { value: chunk, done }: { value?: Uint8Array; done: boolean } = await this.reader.read();

      let chunkArray: Uint8Array | undefined = chunk;
      if (chunk instanceof ArrayBuffer) {
        chunkArray = new Uint8Array(chunk);
      }

      if (done) {
        isDone = true;
        break;
      }

      if (!chunkArray || chunkArray.length === 0) {
        continue;
      }

      const byteCountToCopy = Math.min(byteCount - startingOffset, chunkArray.length);

      const sourceArray = chunkArray.slice(0, byteCountToCopy);
      returnBytes.set(sourceArray, startingOffset);
      startingOffset += byteCountToCopy;

      if (chunkArray.length > byteCountToCopy) {
        this.remainingChunk = chunkArray;
        this.remainingChunkOffset = byteCountToCopy;
      }
    }

    if (startingOffset < byteCount) {
      returnBytes = returnBytes.slice(0, startingOffset);
    }

    if (startingOffset > 0) {
      isDone = false;
    }

    return { value: returnBytes, done: isDone };
  }
}
