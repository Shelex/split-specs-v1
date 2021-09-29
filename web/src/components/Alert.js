import { useState } from 'react';

const ClosingAlert = ({ message }) => {
    const [showAlert, setShowAlert] = useState(true);

    setTimeout(() => {
        setShowAlert(false);
    }, 10000);

    return (
        <div>
            {showAlert ? (
                <div
                    className={`text-white px-6 py-4 border-0 rounded relative mb-4 bg-red-400`}
                >
                    <span className="text-xl inline-block mr-5 align-middle">
                        <i className="fas fa-bell" />
                    </span>
                    <span className="inline-block align-middle mr-8">
                        <b className="capitalize">Error: </b>
                        {message}
                    </span>
                    <button
                        className="absolute bg-transparent text-2xl font-semibold leading-none right-0 top-0 mt-4 mr-6 outline-none focus:outline-none"
                        onClick={() => setShowAlert(false)}
                    >
                        <span>×</span>
                    </button>
                </div>
            ) : null}
        </div>
    );
};

export default function ShowAlert(message = 'Something went wrong :c') {
    const { graphQLErrors, networkError } = message;
    if (networkError) {
        return <ClosingAlert message={networkError.bodyText} />;
    }
    if (graphQLErrors && graphQLErrors.length) {
        return graphQLErrors.map((graphqlError) => {
            return <ClosingAlert message={graphqlError.message} />;
        });
    }

    return <ClosingAlert message={message} />;
}
