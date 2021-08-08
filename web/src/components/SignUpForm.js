import { memo, useCallback, useState } from 'react';
import { useMutation } from '@apollo/client';

import { isLoggedInVar } from '../apollo';
import { SIGN_UP } from '../apollo/mutation';

import Spinner from '../components/Spinner';

const SignUpForm = ({ history }) => {
    const [signUp, { error, loading }] = useMutation(SIGN_UP, {
        onCompleted: (data) => {
            localStorage.setItem('token', data.register);
            isLoggedInVar(true);
            history.replace('/');
        }
    });

    const [values, setValues] = useState();

    const validate = (values) =>
        !values ||
        !values?.email ||
        !values?.password ||
        !values?.passwordConfirm ||
        values?.password?.length < 4 ||
        values?.password !== values?.passwordConfirm;

    const onSubmit = useCallback(
        (e) => {
            e.preventDefault();

            signUp({
                variables: {
                    ...values
                }
            });
        },
        [signUp, values]
    );

    const onChange = useCallback((e) => {
        setValues((prev) => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
    }, []);

    return (
        <form onSubmit={onSubmit}>
            <div className="max-w-lg mx-auto bg-white p-4">
                <div className="mb-4">
                    <label className="signup-label" htmlFor="email">
                        Email
                    </label>
                    <input
                        className="signup-input"
                        type="email"
                        name="email"
                        id="email"
                        placeholder="Please enter your email"
                        autoComplete="off"
                        required
                        onChange={onChange}
                    />
                </div>

                <div className="mb-4">
                    <label className="signup-label" htmlFor="password">
                        Password
                    </label>
                    <input
                        className="signup-input"
                        type="password"
                        name="password"
                        id="password"
                        placeholder="Please enter your password, min 4 chars"
                        autoComplete="off"
                        required
                        onChange={onChange}
                    />
                </div>

                <div className="mb-4">
                    <label className="signup-label" htmlFor="password-check">
                        Confirm Password
                    </label>
                    <input
                        type="password"
                        name="passwordConfirm"
                        id="password-check"
                        className="signup-input"
                        placeholder="Please enter your password again"
                        autoComplete="off"
                        required
                        onChange={onChange}
                    />
                </div>

                <div className="text-xs font-semibold text-red-500">
                    {error && `${error}`}
                </div>

                <div className="mt-12">
                    <button
                        disabled={validate(values)}
                        type="submit"
                        className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none disabled:opacity-50"
                    >
                        {loading ? <Spinner /> : `Sign up`}
                    </button>
                </div>
            </div>
        </form>
    );
};

export default memo(SignUpForm);
