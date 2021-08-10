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
        <div className="max-w-2xl px-4 mx-auto mt-8">
            <div className="text-2xl">{project?.projectName}</div>
            <div>
                {project?.sessions &&
                    Sessions(project?.projectName, project.sessions)}
            </div>
            <br />
            <br />
            <br />
            <DeleteButton
                title="Delete project"
                onClick={onDelete}
                data={deleteData}
                loading={deleteLoading}
            />
        </div>
    );
};

const Sessions = (project, sessions) => {
    const orderedSessions = [...sessions].sort((a, b) => b?.end - a?.end);
    return (
        <div>
            <p>
                have {orderedSessions?.length}
                {pluralize(' session', orderedSessions?.length)}:
            </p>
            {orderedSessions?.length &&
                orderedSessions.map((sessions) => Session(project, sessions))}
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
        new Set(session.backlog.map((item) => item.assignedTo))
    );

    const expectedSerialDuration = session.backlog
        .map((spec) => spec.estimatedDuration)
        .reduce((a, b) => a + b, 0);

    const savedDuration = secondsToDuration(
        expectedSerialDuration - executionTime
    );

    return (
        <li key={session?.id}>
            <Link
                to={{
                    pathname: `/session/${session.id}`,
                    state: {
                        id: session.id,
                        projectName: project
                    }
                }}
            >
                {session.backlog.length} specs {duration}
                {completed && ` at ${displayStart} - ${displayEnd} `}
                with {sessionMachines.length}
                {pluralize('machine', sessionMachines.length)}
                {completed ? ` saved ${savedDuration}` : ''}
            </Link>
        </li>
    );
};

export default memo(Project);
