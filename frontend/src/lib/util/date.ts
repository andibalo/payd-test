import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";

dayjs.extend(utc);

export function utcToLocal(utcString: string, format = "YYYY-MM-DD HH:mm") {
  return dayjs.utc(utcString).local().format(format);
}