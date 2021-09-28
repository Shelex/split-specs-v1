import { memo, useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import DatePicker from 'react-datepicker';
import dayjs from 'dayjs';
import { useMutation } from '@apollo/client';
import { CREATE_API_KEY } from '../apollo/mutation';

import Spinner from '../components/Spinner';

const CreateApiKeyForm = () => {
    let [createApiKey, { error, data, loading }] = useMutation(CREATE_API_KEY);

    const currentDate = new Date();

    const defaultDate = dayjs().add(3, 'months').toDate();

    const [values, setValues] = useState({
        expireAt: dayjs(defaultDate).unix() * 1000
    });

    const validate = (values) => {
        return !values || !values?.name || !values?.expireAt;
    };

    const onSubmit = useCallback(
        (e) => {
            e.preventDefault();

            const { name, expireAt } = values;

            createApiKey({
                variables: {
                    name: name,
                    expireAt: expireAt.valueOf() / 1000
                }
            });

            values.name = null;
        },
        [createApiKey, values]
    );

    const onChange = useCallback((e) => {
        setValues((prev) => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
    }, []);

    return (
        <div className="min-w-full flex items-center justify-center px-4">
            <div className="max-w-md w-full bg-white rounded-md p-6 shadow-2xl">
                <div className="mb-6">
                    <form onSubmit={onSubmit}>
                        <p>Create new API Key</p>
                        <div className="mx-auto bg-white mt-4">
                            <div className="mb-6">
                                <label className="form-label" htmlFor="name">
                                    Please enter api key name
                                </label>
                                <input
                                    className="form-input"
                                    type="text"
                                    name="name"
                                    placeholder="Please enter name of api key"
                                    autoComplete="on"
                                    required
                                    onChange={onChange}
                                />
                            </div>

                            <div className="mb-6">
                                <label
                                    className="form-label"
                                    htmlFor="specFiles"
                                >
                                    Please enter api key expiry date
                                </label>
                                <DatePicker
                                    name="expireAt"
                                    selected={values?.expireAt || currentDate}
                                    onChange={(date) =>
                                        onChange({
                                            target: {
                                                name: 'expireAt',
                                                value: date
                                            }
                                        })
                                    }
                                    startDate={defaultDate}
                                    minDate={new Date()}
                                    nextMonthButtonLabel=">"
                                    previousMonthButtonLabel="<"
                                    dateFormat="yyyy-MM-dd"
                                />
                            </div>

                            <div className="text-xs font-semibold text-red-500">
                                {error}
                            </div>

                            {data?.addApiKey && (
                                <div>
                                    <div className="text-xs font-semibold my-2">
                                        {' '}
                                        API Key created:
                                    </div>
                                    <div className="text-xs font-semibold text-green-500 break-all">
                                        {data?.addApiKey}
                                    </div>
                                    <button className="bg-blue-500 hover:bg-blue-700 my-2 text-white font-bold py-1 px-2 rounded">
                                        <Link to={`/apikeys`}>
                                            Back to api keys list
                                        </Link>
                                    </button>
                                </div>
                            )}

                            <div className="mt-12">
                                <button
                                    disabled={validate(values)}
                                    type="submit"
                                    className="bg-blue-800 w-full py-3 rounded-md text-white hover:bg-blue-900 focus:outline-none disabled:opacity-50"
                                >
                                    {loading ? <Spinner /> : `Create Api Key`}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default memo(CreateApiKeyForm);
