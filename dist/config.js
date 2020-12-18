"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const bunyan_1 = __importDefault(require("bunyan"));
exports.default = {
    logger: bunyan_1.default.createLogger({
        name: 'orchid',
        level: process.env.ORCHID_LOG_LEVEL || 'debug',
    })
};
