import Turbolinks from "turbolinks";
import { Application } from "@hotwired/stimulus";
import axios from "axios";
import { init } from "./core";
import { restoreForm } from "./form";
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
import {
  FormatDateRelativeController,
  FormatInputDateController,
} from "./date";
import { ClickLinkSubmitFormController, XHRSubmitFormController } from "./form";
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

Stimulus.register("format-date-relative", FormatDateRelativeController);
Stimulus.register("format-input-date", FormatInputDateController);

Stimulus.register("click-link-submit-form", ClickLinkSubmitFormController);
Stimulus.register("xhr-submit-form", XHRSubmitFormController);

window.api.onLoad(() => {
  document.body.classList.add("js");
});

window.api.onLoad(setupModal);

window.api.onLoad(restoreForm);
