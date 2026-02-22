import { describe, it, expect } from "vitest";
import {
  convertSmallBufferToString,
  convertSmallStringToBuffer,
  convertSmallUint8ArrayToString,
} from "./buffer-utils.js";

describe("buffer-utils", () => {
  describe("convertSmallBufferToString / convertSmallStringToBuffer", () => {
    it("round-trips ArrayBuffer", () => {
      const original = new Uint8Array([72, 101, 108, 108, 111]).buffer;
      const encoded = convertSmallBufferToString(original);
      const decoded = convertSmallStringToBuffer(encoded);
      expect(new Uint8Array(decoded)).toEqual(new Uint8Array(original));
    });

    it("handles empty buffer", () => {
      const original = new ArrayBuffer(0);
      const encoded = convertSmallBufferToString(original);
      const decoded = convertSmallStringToBuffer(encoded);
      expect(decoded.byteLength).toBe(0);
    });

    it("handles binary data", () => {
      const original = new Uint8Array([0, 255, 128, 1]);
      const encoded = convertSmallBufferToString(original.buffer);
      const decoded = convertSmallStringToBuffer(encoded);
      expect(new Uint8Array(decoded)).toEqual(original);
    });
  });

  describe("convertSmallUint8ArrayToString", () => {
    it("encodes Uint8Array to base64 string", () => {
      const arr = new Uint8Array([72, 101, 108, 108, 111]);
      const encoded = convertSmallUint8ArrayToString(arr);
      expect(encoded).toBe("SGVsbG8=");
    });
  });
});
