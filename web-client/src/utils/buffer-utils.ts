export const convertSmallBufferToString = (buffer: ArrayBuffer): string => {
  return window.btoa(String.fromCharCode.apply(null, Array.from(new Uint8Array(buffer))));
};

export const convertSmallUint8ArrayToString = (array: Uint8Array): string => {
  return window.btoa(String.fromCharCode.apply(null, Array.from(array)));
};

export const convertSmallStringToBuffer = (packed: string): ArrayBuffer => {
  const string = window.atob(packed);
  const buffer = new ArrayBuffer(string.length);
  const bufferView = new Uint8Array(buffer);

  for (let i = 0; i < string.length; i++) {
    bufferView[i] = string.charCodeAt(i);
  }

  return buffer;
};
