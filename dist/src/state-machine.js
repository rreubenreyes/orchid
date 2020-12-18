"use strict";
/* This project is meant to output JSONs which are valid Amazon States Language files.
 *
 * See: https://states-language.net/spec.html
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.SFNStateKind = void 0;
var SFNStateKind;
(function (SFNStateKind) {
    SFNStateKind["Task"] = "Task";
    SFNStateKind["Parallel"] = "Parallel";
    SFNStateKind["Map"] = "Map";
    SFNStateKind["Pass"] = "Pass";
    SFNStateKind["Wait"] = "Wait";
    SFNStateKind["Choice"] = "Choice";
    SFNStateKind["Succeed"] = "Succeed";
    SFNStateKind["Fail"] = "Fail";
})(SFNStateKind = exports.SFNStateKind || (exports.SFNStateKind = {}));
