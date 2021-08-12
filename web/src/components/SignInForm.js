import { memo, useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import { useMutation } from '@apollo/client';

import { isLoggedInVar } from '../apollo';
import { SIGN_IN } from '../apollo/mutation';

import Spinner from '../components/Spinner';

const SignInForm = () => {
    const [signIn, { error, loading }] = useMutation(SIGN_IN, {
        onCompleted: (data) => {
            localStorage.setItem('token', data.login);
            isLoggedInVar(true);
        }
    });

    const [values, setValues] = useState();

    const onSubmit = useCallback(
        (e) => {
            e.preventDefault();

            signIn({
                variables: {
                    ...values
                }
            });
        },
        [signIn, values]
    );

    const onChange = useCallback((e) => {
        setValues((prev) => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
    }, []);

    return (
        <div className="fixed inset-0 bg-gray-900 bg-opacity-60">
            <div className="min-w-full min-h-full flex items-center justify-center px-4">
                <div className="max-w-md w-full bg-white rounded-md p-6 shadow-2xl">
                    <div className="mb-6">
                        <h2 className="text-center text-3xl font-extrabold text-gray-700 mt-4">
                            Split specs
                        </h2>
                        <p className="text-center text-xs text-gray-600 mt-1">
                            <Link
                                className="text-blue-600 font-semibold"
                                to="/signup"
                            >
                                I am new here
                            </Link>
                        </p>
                    </div>

                    <form onSubmit={onSubmit}>
                        <div className="max-w-lg mx-auto mb-2">
                            <div>
                                <input
                                    className="form-input"
                                    type="email"
                                    name="email"
                                    placeholder="Please enter your email"
                                    autoComplete="off"
                                    required
                                    onChange={onChange}
                                />
                            </div>

                            <div>
                                <input
                                    className="form-input"
                                    type="password"
                                    name="password"
                                    placeholder="Please enter your password"
                                    autoComplete="off"
                                    required
                                    onChange={onChange}
                                />
                            </div>
                        </div>

                        <div className="text-xs font-semibold text-red-500">
                            {error && `${error.message}`}
                        </div>

                        <div className="mt-8">
                            <button
                                type="submit"
                                className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none"
                            >
                                {loading ? <Spinner /> : <p>Sign in</p>}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default memo(SignInForm);
