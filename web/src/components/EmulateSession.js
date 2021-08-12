import { useLazyQuery, useQuery } from '@apollo/client';
import { memo, useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import { timestampToDate, secondsToDuration } from '../format/displayDate';

import { GET_SESSION, NEXT_SPEC } from '../apollo/query';

import Spinner from './Spinner';

export const EmulateSession = ({ session }) => {
    const { sessionId, projectName } = session;

    const { data, refetch: refetchSession } = useQuery(GET_SESSION, {
        variables: {
            id: sessionId
        },
        fetchPolicy: 'network-only',
        nextFetchPolicy: 'cache-and-network'
    });

    const [values, setValues] = useState();

    const onChange = useCallback((e) => {
        setValues((prev) => ({
            ...prev,
            [e.target.name]: e.target.value
        }));
    }, []);

    const [getNextSpec, { loading: nextSpecLoading, error: nextSpecError }] =
        useLazyQuery(NEXT_SPEC, {
            fetchPolicy: 'no-cache',
            nextFetchPolicy: 'no-cache',
            errorPolicy: 'ignore'
        });

    const onNextSpec = (machineId) => async (e) => {
        e.preventDefault();
        const options = {
            variables: {
                sessionId,
                options: {
                    machineId: machineId || 'default'
                }
            }
        };
        getNextSpec(options);

        setValues((prev) => ({
            ...prev,
            nextSpec: values?.nextSpec
        }));

        await new Promise((resolve) =>
            setTimeout(() => {
                refetchSession({
                    variables: {
                        id: sessionId
                    }
                });
                resolve();
            }, 500)
        );
    };

    return (
        <div className="max-w-6xl px-4 mx-auto">
            <p>Session created</p>
            <p>project: {projectName}</p>
            <p>id: {sessionId}</p>
            <p className="w-max">
                <Link
                    to={`/session/${projectName}/${sessionId}`}
                    location={projectName}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    <button className="bg-green-500 w-full px-2  py-3 rounded-md text-white hover:bg-green-700 focus:outline-none disabled:opacity-50">
                        open session
                    </button>
                </Link>
            </p>
            <div>
                <SpecsTable
                    specs={data?.session?.backlog}
                    finished={data?.session?.end > 0}
                    started={data?.session?.start > 0}
                />
            </div>
            <div className="mt-5">
                <input
                    className="form-input"
                    type="text"
                    name="machineId"
                    defaultValue={values?.machineId || 'default'}
                    placeholder="Please enter name of current machine"
                    autoComplete="on"
                    required
                    onChange={onChange}
                />
                <div className="text-xs font-semibold text-red-500">
                    {nextSpecError && `${nextSpecError}`}
                </div>
                <button
                    className={`bg-green-500 hover:bg-green-700 text-white font-bold py-3 px-2 mt-5 rounded w-full`}
                    onClick={onNextSpec(values?.machineId)}
                >
                    {nextSpecLoading ? (
                        <Spinner />
                    ) : (
                        <p>Request next spec for {values?.machineId}</p>
                    )}
                </button>
            </div>
        </div>
    );
};

const SpecsTable = ({ specs, started, finished }) => {
    const executionTime = 'Execution time';
    const estimatedDuration = 'Estimated Duration';

    const durationHeader =
        started && finished
            ? executionTime
            : started
            ? `${estimatedDuration} / ${executionTime}`
            : estimatedDuration;

    return specs ? (
        <div className="mt-5">
            <table className="table-auto border-collapse border border-blue-400">
                <thead className="space-x-1">
                    <tr className="bg-blue-600 px-auto py-auto">
                        <th className="w-1/3">
                            <span className="text-gray-100 font-semibold">
                                Name
                            </span>
                        </th>
                        <th className="w-1/5">
                            <span className="text-gray-100 font-semibold">
                                {durationHeader}
                            </span>
                        </th>

                        <th className="w-1/4">
                            <span className="text-gray-100 font-semibold">
                                Start
                            </span>
                        </th>

                        <th className="w-1/4">
                            <span className="text-gray-100 font-semibold">
                                End
                            </span>
                        </th>

                        <th className="w-1/2">
                            <span className="text-gray-100 font-semibold">
                                Machine
                            </span>
                        </th>
                    </tr>
                </thead>
                <tbody className="bg-gray-200">
                    {specs.map((spec) => (
                        <tr key={spec.file} className="bg-white">
                            <td className="font-semibold border border-blue-400">
                                {spec.file}
                            </td>
                            <td className="border border-blue-400">
                                {secondsToDuration(spec.estimatedDuration)}
                            </td>

                            <td className="border border-blue-400">
                                {spec.start > 0
                                    ? timestampToDate(spec.start)
                                    : ''}
                            </td>
                            <td className="border border-blue-400">
                                {spec.end > 0 ? timestampToDate(spec.end) : ''}
                            </td>

                            <td className="border border-blue-400">
                                {spec.assignedTo || 'none'}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    ) : (
        <p>no specs received</p>
    );
};

export default memo(EmulateSession);
