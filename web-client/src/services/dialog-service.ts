import { Dialog, Notify } from "quasar";

export const NotificationType = {
  SUCCESS: "SUCCESS",
  ERROR: "ERROR",
  INFO: "INFO",
  WARNING: "WARNING",
} as const;

export const dialogService = {
  alert(title: string, message: string): Promise<boolean> {
    return new Promise((accept) => {
      Dialog.create({
        title,
        message,
      })
        .onOk(() => accept(true))
        .onDismiss(() => accept(false));
    });
  },

  confirm(title: string, message: string): Promise<boolean> {
    return new Promise((accept) => {
      Dialog.create({
        title,
        message,
        cancel: true,
        persistent: true,
      })
        .onOk(() => accept(true))
        .onCancel(() => accept(false))
        .onDismiss(() => accept(false));
    });
  },

  prompt(title: string, message: string, initialValue = ""): Promise<string | null> {
    return new Promise((accept) => {
      Dialog.create({
        title,
        message,
        prompt: {
          model: initialValue,
          type: "text",
        },
        cancel: true,
        persistent: true,
      })
        .onOk((answer: string) => accept(answer))
        .onCancel(() => accept(null))
        .onDismiss(() => accept(null));
    });
  },

  notify(type: string, message: string): void {
    let color = "positive";
    let icon = "info";

    if (type === NotificationType.SUCCESS) {
      color = "positive";
      icon = "check_circle";
    } else if (type === NotificationType.ERROR) {
      color = "negative";
      icon = "error";
    } else if (type === NotificationType.WARNING) {
      color = "warning";
      icon = "warning";
    } else {
      color = "info";
      icon = "info";
    }

    Notify.create({
      message,
      color,
      icon,
      position: "top",
      timeout: 3000,
    });
  },
};
