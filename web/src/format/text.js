export const pluralize = (message, count) => {
    const shouldBePlural = count % 10 !== 1;
    return shouldBePlural ? `${message}s` : message;
};

export const capitalize = (word) => {
    return word && word.charAt(0).toUpperCase() + word.slice(1);
};
