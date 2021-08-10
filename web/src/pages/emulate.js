import { memo } from 'react';

import SessionForm from '../components/CreateSessionForm';

const Emulate = () => {
    return (
        <div className="h-full pt-4 sm:pt-12">
            <div className="max-w-lg mx-auto p-4">
                <div>
                    <h2 className="text-xl font-semibold text-gray-700">
                        Emulate session
                    </h2>
                </div>
            </div>

            <SessionForm history={history} />
        </div>
    );
};

export default memo(Emulate);
