import { memo, useCallback } from 'react';
import { useMutation, useQuery } from '@apollo/client';
import { displayTimestamp, secondsToDuration } from '../format/displayDate';
import Loading from '../components/Loading';
import { DeleteButton } from '../components/DeleteButton';

import { GET_SESSION } from '../apollo/query';
import { DELETE_SESSION } from '../apollo/mutation';
import { Link, useHistory, useParams } from 'react-router-dom';

const Session = () => {
    const { projectName, id } = useParams();

    const [deleteSession, { data: deleteData, loading: deleteLoading }] =
        useMutation(DELETE_SESSION);

    const { data, loading } = useQuery(GET_SESSION, {
        variables: { id },
        fetchPolicy: 'network-only'
    });

    const history = useHistory();

    const onDelete = useCallback(
        (e) => {
            e.preventDefault();
            deleteSession({
                variables: { sessionId: id }
            }).then(() => {
                history.push(`/project/${projectName}`);
            });
        },
        [deleteSession]
    );

    const session = data?.session;

    if (loading) {
        return <Loading />;
    }

    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            {projectName && (
                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded">
                    <Link to={`/project/${projectName}`}>Back to project</Link>
                </button>
            )}
            <div className="text-2xl">Session id "{session?.id}"</div>
            <div className="text-1xl">
                {displayTimestamp(session?.start)} :{' '}
                {displayTimestamp(session?.end)}
            </div>
            {Specs(
                session?.backlog,
                session.start > 0 && session.end > 0,
                projectName
            )}
            <div className="mt-10">{ByMachine(session)}</div>
            <DeleteButton
                title="Delete session"
                onClick={onDelete}
                data={deleteData}
                loading={deleteLoading}
            />
        </div>
    );
};

const Specs = (backlog, completed, projectName) => {
    return (
        <div>
            <div>{backlog.length} files</div>
            {backlog.length > 0 && (
                <table className="table-auto border-collapse border border-blue-400 w-full">
                    <thead className="space-x-1">
                        <tr className="bg-blue-600 px-auto py-auto">
                            <th className="w-1/2 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    FileName
                                </span>
                            </th>
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    {completed
                                        ? 'Duration'
                                        : 'Estimated duration'}
                                </span>
                            </th>
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    Machine
                                </span>
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-gray-200">
                        {backlog?.length &&
                            [...backlog]
                                .sort(
                                    (a, b) =>
                                        b.estimatedDuration -
                                        a.estimatedDuration
                                )
                                .map((spec) => Spec(spec, projectName))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

const Spec = (spec, projectName) => {
    return (
        <tr key={spec.file} className="bg-white">
            <td className="font-semibold border border-blue-400">
                <Link
                    to={`/spec/${projectName}/${encodeURIComponent(spec.file)}`}
                >
                    {spec.file}
                </Link>
            </td>
            <td className="border border-blue-400">
                {secondsToDuration(spec.estimatedDuration)}
            </td>
            <td className="border border-blue-400">{spec.assignedTo}</td>
        </tr>
    );
};

const ByMachine = (session) => {
    const sessionMachines = Array.from(
        new Set(
            session?.backlog.map((item) => item?.assignedTo).filter((x) => x)
        )
    );

    const statsPerMachine = sessionMachines.map((machine) => {
        return {
            machine: machine,
            duration: session?.backlog
                .filter((spec) => spec.assignedTo === machine)
                .map((spec) => spec.estimatedDuration)
                .reduce((a, b) => a + b, 0)
        };
    });

    return (
        <div>
            <div>{sessionMachines.length} machines</div>
            {sessionMachines.length > 0 && (
                <table className="table-auto border-collapse border border-blue-400 w-full">
                    <thead className="space-x-1">
                        <tr className="bg-blue-600 px-auto py-auto">
                            <th className="w-1/2 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    MachineID
                                </span>
                            </th>
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    Duration
                                </span>
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-gray-200">
                        {statsPerMachine
                            .sort((a, b) => a.machine.localeCompare(b.machine))
                            .map((stat) => (
                                <tr key={stat.machine} className="bg-white">
                                    <td className="font-semibold border border-blue-400">
                                        {stat.machine}
                                    </td>
                                    <td className="border border-blue-400">
                                        {secondsToDuration(stat.duration)}
                                    </td>
                                </tr>
                            ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default memo(Session);
