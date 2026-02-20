export const ClipboardAction = {
  CUT: "CUT",
  COPY: "COPY",
} as const;

export type ClipboardActionType = (typeof ClipboardAction)[keyof typeof ClipboardAction];
