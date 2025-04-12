export const SECONDS_PER_MINUTE = 60;
export const SECONDS_PER_HOUR = 60 * SECONDS_PER_MINUTE;
export const SECONDS_PER_DAY = 24 * SECONDS_PER_HOUR;
export const SECONDS_PER_WEEK = 7 * SECONDS_PER_DAY;

export const formatDurationSeconds = (durationSeconds: number, delimeter?: string) => {
    let remainingSeconds = durationSeconds;
    const parts = [];

    const weeks = Math.floor(remainingSeconds / SECONDS_PER_WEEK);
    if (weeks) {
        parts.push(`${weeks}w`);
    }
    remainingSeconds -= weeks * SECONDS_PER_WEEK;

    const days = Math.floor(remainingSeconds / SECONDS_PER_DAY);
    if (days) {
        parts.push(`${days}d`);
    }
    remainingSeconds -= days * SECONDS_PER_DAY;

    const hours = Math.floor(remainingSeconds / SECONDS_PER_HOUR);
    if (hours) {
        parts.push(`${hours}h`);
    }
    remainingSeconds -= hours * SECONDS_PER_HOUR;

    const minutes = Math.floor(remainingSeconds / SECONDS_PER_MINUTE);
    if (minutes) {
        parts.push(`${minutes}m`);
    }
    remainingSeconds -= minutes * SECONDS_PER_MINUTE;

    if (remainingSeconds) {
        parts.push(`${remainingSeconds}s`);
    }
    return parts.join(delimeter || '');
};

export const parseDuration = (duration: string) => {
    const regex = /(?:(\d+)w)?(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?/;
    const match = duration.match(regex);

    if (!match) {
        throw new Error(`Invalid duration format: ${duration}`);
    }

    const weeks = parseInt(match[1] || '0', 10);
    const days = parseInt(match[2] || '0', 10);
    const hours = parseInt(match[3] || '0', 10);
    const minutes = parseInt(match[4] || '0', 10);
    const seconds = parseInt(match[5] || '0', 10);

    return (
        weeks * SECONDS_PER_WEEK +
        days * SECONDS_PER_DAY +
        hours * SECONDS_PER_HOUR +
        minutes * SECONDS_PER_MINUTE +
        seconds
    );
};

export const formatTimeRange = (time?: Date, durationSeconds?: number) => {
    if (!time) {
        return 'All Time';
    } else if (!durationSeconds) {
        return `Since ${time.toLocaleString()}`;
    }
    return `${time.toLocaleString()} + ${formatDurationSeconds(durationSeconds, ' ')}`;
};
