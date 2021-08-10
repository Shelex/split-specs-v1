import { memo, useCallback } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { Link, useParams } from 'react-router-dom';
import Loading from '../components/Loading';
import { displayTimestamp, secondsToDuration } from '../format/displayDate';
import { pluralize } from '../format/text';

import { DeleteButton } from '../components/DeleteButton';
import { GET_PROJECT } from '../apollo/query';
import { DELETE_PROJECT } from '../apollo/mutation';

const Project = () => {
    const { name } = useParams();
    const { data, loading } = useQuery(GET_PROJECT, {
        variables: { name }
    });

    const [deleteProject, { data: deleteData, loading: deleteLoading }] =
        useMutation(DELETE_PROJECT);

    const onDelete = useCallback(
        (e) => {
            e.preventDefault();
            deleteProject({
                variables: {
                    projectName: name
                }
            });
        },
        [deleteProject]
    );

    if (loading) {
        return <Loading />;
    }

    const project = data?.project;

    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            <div className="text-2xl">{project?.projectName}</div>
            <div>
                {project?.sessions &&
                    Sessions(project?.projectName, project.sessions)}
                <DeleteButton
                    title="Delete project"
                    onClick={onDelete}
                    data={deleteData}
                    loading={deleteLoading}
                />
            </div>
        </div>
    );
};

const Sessions = (project, sessions) => {
    const orderedSessions = [...sessions].sort((a, b) => b?.end - a?.end);
    return (
        <div>
            <p>
                {orderedSessions?.length}
                {pluralize(' session', orderedSessions?.length)}
            </p>

            {orderedSessions?.length > 0 && (
                <div>
                    <table className="table-auto border-collapse border border-blue-400">
                        <thead className="space-x-1">
                            <tr className="bg-blue-600 px-auto py-auto">
                                <th className="w-1/5 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        ID
                                    </span>
                                </th>
                                <th className="w-1/8 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Spec amount
                                    </span>
                                </th>
                                <th className="w-1/8 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Duration
                                    </span>
                                </th>

                                <th className="w-1/5 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Start
                                    </span>
                                </th>

                                <th className="w-1/5 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        End
                                    </span>
                                </th>

                                <th className="w-1/8 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Machines
                                    </span>
                                </th>
                                <th className="w-1/8 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Saved time
                                    </span>
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-gray-200">
                            {orderedSessions?.length &&
                                orderedSessions.map((sessions) =>
                                    Session(project, sessions)
                                )}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

const Session = (project, session) => {
    const displayStart = displayTimestamp(session.start);
    const displayEnd = displayTimestamp(session.end);

    const executionTime = session.end - session.start;
    const executionTimeMessage = secondsToDuration(executionTime);

    const isStarted = session.start > 0;
    const isFinished = session.end > 0;

    const completed = isStarted && isFinished;

    const uncompletedDurationMessage = isStarted
        ? 'not finished '
        : 'not started ';

    const duration = completed
        ? executionTimeMessage
        : uncompletedDurationMessage;

    const sessionMachines = Array.from(
        new Set(session.backlog.map((item) => item.assignedTo).filter((x) => x))
    );

    const expectedSerialDuration = session.backlog
        .map((spec) => spec.estimatedDuration)
        .reduce((a, b) => a + b, 0);

    const savedDuration = secondsToDuration(
        expectedSerialDuration - executionTime
    );

    return (
        <tr key={session?.id} className="bg-white">
            <td className="font-semibold border border-blue-400">
                <Link
                    to={{
                        pathname: `/session/${project}/${session.id}`,
                        state: {
                            id: session.id,
                            projectName: project
                        }
                    }}
                >
                    {session?.id}
                </Link>
            </td>
            <td className="border border-blue-400">{session.backlog.length}</td>
            <td className="border border-blue-400">{duration}</td>
            <td className="border border-blue-400">{displayStart}</td>
            <td className="border border-blue-400">{displayEnd}</td>
            <td className="border border-blue-400">{sessionMachines.length}</td>
            <td className="border border-blue-400">
                {completed ? savedDuration : 0}
            </td>
        </tr>
    );
};

export default memo(Project);
