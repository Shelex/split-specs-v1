import { memo, useEffect } from 'react';
import SignUpForm from '../components/SignUpForm';

const SignUp = ({ history }) => {
    useEffect(() => {
        document.title = 'Split Specs: SignUp';
    }, []);

    return (
        <div className="h-full pt-4 sm:pt-12">
            <div className="max-w-lg mx-auto p-4">
                <div>
                    <h2 className="text-xl font-semibold text-gray-700">
                        SignUp
                    </h2>
                    <p className="mt-1 text-gray-600">
                        You will never know until you try
                    </p>
                </div>
            </div>

            <SignUpForm history={history} />
        </div>
    );
};

export default memo(SignUp);
