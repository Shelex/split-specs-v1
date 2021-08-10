import { useCallback } from 'react';
import { useHistory } from 'react-router-dom';

import Spinner from './Spinner';

export const DeleteButton = ({ onClick, loading, data, title }) => {
    let history = useHistory();
    const GoToPreviousPath = () => (
        <div>
            {title.includes('session') ? history.goBack() : history.push('/')}
        </div>
    );

    const onConfirm = useCallback((e) => {
        e.preventDefault();
        if (confirm('Are you really sure?')) {
            return onClick(e);
        }
    });

    return (
        <div className="mt-10">
            {data ? (
                <GoToPreviousPath />
            ) : (
                <button
                    className={`bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-2 rounded w-48`}
                    onClick={onConfirm}
                >
                    {loading ? <Spinner /> : <p>{title}</p>}
                </button>
            )}
        </div>
    );
};
