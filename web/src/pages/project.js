import { memo } from 'react';
import { useQuery } from '@apollo/client';
import { Link, useParams } from 'react-router-dom';
import Loading from '../components/Loading';
import { displayTimestamp, secondsToDuration } from '../dates/displayDate';

import { GET_PROJECT } from '../apollo/query';

const Project = () => {
    const { name } = useParams();
    const { data, loading } = useQuery(GET_PROJECT, {
        variables: { name }
    });

    if (loading) {
        return <Loading />;
    }

    const project = data?.project;

    return (
        <div className="max-w-7xl px-4 mx-auto mt-8">
            <div className="text-2xl">{project?.projectName}</div>
            <div>
                {project?.sessions &&
                    Sessions(project?.projectName, project.sessions)}
            </div>
        </div>
    );
};

const Sessions = (project, sessions) => {
    const orderedSessions = [...sessions].sort((a, b) => b?.end - a?.end);
    return (
        <div>
            <p>have {orderedSessions?.length} sessions:</p>
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

    const shouldBePlural = sessionMachines.length % 10 !== 1;

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
                with {sessionMachines.length} machine
                {shouldBePlural ? 's' : ''}
                {completed ? ` saved ${savedDuration}` : ''}
            </Link>
        </li>
    );
};

export default memo(Project);
