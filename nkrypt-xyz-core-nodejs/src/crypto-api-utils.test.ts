import { describe, it, expect } from "vitest";
import { buildCryptoHeader, unbuildCryptoHeader } from "./crypto-api-utils.js";
import { convertSmallStringToBuffer } from "./buffer-utils.js";

describe("crypto-api-utils", () => {
  describe("buildCryptoHeader / unbuildCryptoHeader", () => {
    it("round-trips iv and salt", () => {
      const iv = "dGVzdC1pdi1kYXRh";
      const salt = "dGVzdC1zYWx0LWRhdGE=";
      const header = buildCryptoHeader(iv, salt);
      expect(header).toContain("NK001|");
      expect(header).toContain(iv);
      expect(header).toContain(salt);

      const [decodedIv, decodedSalt] = unbuildCryptoHeader(header);
      expect(new Uint8Array(decodedIv)).toEqual(new Uint8Array(convertSmallStringToBuffer(iv)));
      expect(new Uint8Array(decodedSalt)).toEqual(new Uint8Array(convertSmallStringToBuffer(salt)));
    });

    it("splits header correctly", () => {
      const iv = "aaa";
      const salt = "bbb";
      const header = buildCryptoHeader(iv, salt);
      const parts = header.split("|");
      expect(parts[0]).toBe("NK001");
      expect(parts[1]).toBe(iv);
      expect(parts[2]).toBe(salt);
    });
  });
});
