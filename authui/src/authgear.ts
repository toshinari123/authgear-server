import Turbolinks from "turbolinks";
import { Application } from "@hotwired/stimulus";
import axios from "axios";
import { init } from "./core";
import {
  clickLinkSubmitForm,
  autoSubmitForm,
  xhrSubmitForm,
  restoreForm,
} from "./form";
import { formatDateRelative } from "./date";
import { setupModal } from "./modal";
import { CopyButtonController } from "./copy";
import { PasswordVisibilityToggleController } from "./passwordVisibility";
import { PasswordPolicyController } from "./password-policy";
import { ClickToSwitchController } from "./clickToSwitch";
import { PreventDoubleTapController } from "./preventDoubleTap";
import { AccountDelectionController } from "./accountdeletion";
import { ResendButtonController } from "./resendButton";
import { MessageBarController } from "./messageBar";
import { IntlTelInputController } from "./intlTelInput";
import { SelectEmptyValueController, GenderSelectController } from "./select";
import { ImagePickerController } from "./imagepicker";
import { WebSocketController } from "./websocket";
import { FormatInputDateController } from "./date";
// FIXME(css): Build CSS files one by one with another tool
// webpack bundles all CSS files into one bundle.

axios.defaults.withCredentials = true;

init();

const Stimulus = Application.start();
Stimulus.register(
  "password-visibility-toggle",
  PasswordVisibilityToggleController
);
Stimulus.register("password-policy", PasswordPolicyController);
Stimulus.register("click-to-switch", ClickToSwitchController);

Stimulus.register("copy-button", CopyButtonController);

Stimulus.register("prevent-double-tap", PreventDoubleTapController);

Stimulus.register("account-delection", AccountDelectionController);

Stimulus.register("resend-button", ResendButtonController);

Stimulus.register("message-bar", MessageBarController);

Stimulus.register("intl-tel-input", IntlTelInputController);

Stimulus.register("select-empty-value", SelectEmptyValueController);
Stimulus.register("gender-select", GenderSelectController);

Stimulus.register("image-picker", ImagePickerController);

Stimulus.register("websocket", WebSocketController);

Stimulus.register("format-input-date", FormatInputDateController);

window.api.onLoad(() => {
  document.body.classList.add("js");
});

window.api.onLoad(formatDateRelative);

window.api.onLoad(setupModal);

window.api.onLoad(autoSubmitForm);
window.api.onLoad(clickLinkSubmitForm);
window.api.onLoad(xhrSubmitForm);
window.api.onLoad(restoreForm);
