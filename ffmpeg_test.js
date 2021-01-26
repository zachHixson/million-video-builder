const DIGITS = 50;

let clipString = "";
let filterString = "";

for (let i = 0; i < DIGITS; i++){
    let pow = Math.pow(10, i);
    let digit = Math.floor(Math.PI * pow) % 10;

    clipString += "-i clips\\\\" + digit + ".mp4 "
    filterString += "[" + i + ":v][" + i + ":a] "
}

console.log("ffmpeg " + clipString + "-filter_complex \"" + filterString + "concat=n=" + DIGITS + ":v=1:a=1 [v] [a]\" -map \"[v]\" -map \"[a]\" output.mp4")