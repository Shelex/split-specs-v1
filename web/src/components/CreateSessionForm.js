import { memo, useCallback, useState } from 'react';
import { useMutation } from '@apollo/client';

import { CREATE_SESSION } from '../apollo/mutation';

import Spinner from '../components/Spinner';
import Loading from '../components/Loading';
import EmulateSession from './EmulateSession';

const CreateSessionForm = () => {
    const [createSession, { error, data: sessionData, loading }] =
        useMutation(CREATE_SESSION);

    const [values, setValues] = useState();

    const validate = (values) => !values || !values?.projectName;

    const defaultSpecs = 'a,b,c,d,e,f';

    const onSubmit = useCallback(
        (e) => {
            e.preventDefault();

            const { projectName, specFiles } = values;

            const files = (specFiles || defaultSpecs)
                .split(',')
                .filter((x) => x)
                .map((fileName) => ({
                    filePath: fileName.trim()
                }));

            createSession({
                variables: {
                    session: {
                        projectName: projectName,
                        specFiles: files
                    }
                }
            });
        },
        [createSession, values]
    );

    const onChange = useCallback((e) => {
        setValues((prev) => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
    }, []);

    if (loading) {
        return <Loading />;
    }

    return sessionData ? (
        <EmulateSession session={sessionData?.addSession} />
    ) : (
        <div className="min-w-full flex items-center justify-center px-4">
            <div className="max-w-md w-full bg-white rounded-md p-6 shadow-2xl">
                <div className="mb-6">
                    <form onSubmit={onSubmit}>
                        <p>Emulate new session</p>
                        <div className="mx-auto bg-white mt-4">
                            <div className="mb-6">
                                <label
                                    className="signup-label"
                                    htmlFor="projectName"
                                >
                                    Please enter project name
                                </label>
                                <input
                                    className="signup-input"
                                    type="text"
                                    name="projectName"
                                    placeholder="Please enter name of project"
                                    autoComplete="on"
                                    required
                                    onChange={onChange}
                                />
                            </div>

                            <div className="mb-6">
                                <label
                                    className="signup-label"
                                    htmlFor="specFiles"
                                >
                                    Please enter comma-separated spec files
                                </label>
                                <input
                                    className="signup-input"
                                    type="text"
                                    name="specFiles"
                                    placeholder="Please enter spec files"
                                    autoComplete="off"
                                    defaultValue={defaultSpecs}
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
                                    {loading ? <Spinner /> : `Create Session`}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default memo(CreateSessionForm);
