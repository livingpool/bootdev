import { moo } from "./moo.js";
import * as cowsay from "cowsay";

const result = moo("there");
console.log(cowsay.say({ text: result }));
