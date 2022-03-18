export const defineSpecStatusTextAndColor = (spec) => {
    if (!spec.start) {
        return ['', 'white'];
    }
    if (!spec.end) {
        return ['running', 'yellow-300'];
    }
    return spec.passed ? ['passed', 'green-600'] : ['failed', 'red-600'];
};
