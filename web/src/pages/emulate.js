import { memo } from 'react';

import SessionForm from '../components/CreateSessionForm';

const Emulate = () => {
    return (
        <div className="h-full pt-4 sm:pt-12">
            <SessionForm history={history} />
        </div>
    );
};

export default memo(Emulate);
