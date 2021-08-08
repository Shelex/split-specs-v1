import { memo } from 'react';
import Spinner from './Spinner';

const Loading = () => {
    return (
        <div>
            <div className="fixed top-0 left-0 right-0 bottom-0 w-full h-screen z-50 flex flex-col items-center justify-center">
                <Spinner />
                <h2 className="text-center text-white text-xl font-semibold">
                    Loading...
                </h2>
                <p className="w-1/3 text-center text-white">
                    This may take a few seconds
                </p>
            </div>
        </div>
    );
};

export default memo(Loading);
