import { useHistory } from 'react-router-dom';

import Spinner from './Spinner';

export const DeleteButton = ({ onClick, loading, data, title }) => {
    let history = useHistory();
    const GoToPreviousPath = () => (
        <div>
            {title.includes('session')
                ? history.goBack()
                : (window.location.href = '/')}
        </div>
    );

    return (
        <div className="mt-10">
            {data ? (
                <GoToPreviousPath />
            ) : (
                <button
                    className={`bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded w-max`}
                    onClick={onClick}
                >
                    {loading ? <Spinner /> : <p>{title}</p>}
                </button>
            )}
        </div>
    );
};
