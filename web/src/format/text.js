export const pluralize = (message, count) => {
    const shouldBePlural = count % 10 !== 1;
    return shouldBePlural ? `${message}s` : message;
};
