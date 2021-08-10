import dayjs from 'dayjs';

export const timestampToDate = (timestamp) => {
    return dayjs(timestamp * 1000).format('DD-MM-YYYY HH:mm:ss');
};

export const displayTimestamp = (timestamp) => {
    return timestamp > 0 ? timestampToDate(timestamp) : '_';
};

export const secondsToDuration = (seconds) => {
    const temporaryExecutionTime = new Date(0);
    temporaryExecutionTime.setSeconds(seconds);
    return temporaryExecutionTime.toISOString().substr(11, 8);
};
