
import { Readable } from "stream";

import {
  ReadableStream
} from 'node:stream/web';

export const isANodejsReadableStream = (object) => {
  try {
    return object instanceof Readable;
  } catch (ex) {
    console.error(ex);
    return false;
  }
};

export const isABrowserReadableStream = (object) => {
  try {
    return object instanceof ReadableStream;
  } catch (ex) {
    console.error(ex);
    return false;
  }
};