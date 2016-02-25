/* @flow */

const byteLabels = ["B", "KiB", "MiB", "GiB", "TiB"];

export function formatBytes(bytes: number): string {
  for (let i = byteLabels.length-1; i>0; i--) {
    if (bytes>=Math.pow(1000, i)) return (bytes/Math.pow(1024, i)).toFixed(2) + " " + byteLabels[i];
  }
  return bytes + " Bytes";
}

export function formatSeconds(seconds: number): string {
  if (seconds==0) return "";
  if (seconds > 60*60*24*7) {
    return "∞"
  } else if (seconds > 60*60*24) {
    var day = Math.floor(seconds / (60*60*24));
    var hour = Math.ceil((seconds % (60*60*24)) / 24);
    var dayStr = day == 1 ? "1 day" : day + " days";
    if (hour == 0) return dayStr;
    else {
      var hourStr = hour == 1 ? "1 hour" : hour + " hours";
      return dayStr + " " + hourStr;
    }
  } else if (seconds > 60*60) {
    var hour = Math.floor(seconds / (60*60));
    var mins = Math.ceil((seconds % (60*60)) / 60);
    var hourStr = hour == 1 ? "1 hour" : hour + " hours";
    if (mins == 0) return hourStr;
    else {
      var minStr = mins == 1 ? "1 minute" : mins + " minutes";
      return hourStr + " " + minStr;
    }
  } else if (seconds > 100) return Math.ceil(seconds/60) + " minutes";
  else return seconds == 1 ? "1 second" : (seconds|0) + " seconds";
}
