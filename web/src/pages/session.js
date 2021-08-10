import { memo, useCallback } from 'react';
import { useMutation, useQuery } from '@apollo/client';
import { displayTimestamp, secondsToDuration } from '../format/displayDate';
import Loading from '../components/Loading';
import { DeleteButton } from '../components/DeleteButton';
import { Redirect } from 'react-router';

import { GET_SESSION } from '../apollo/query';
import { DELETE_SESSION } from '../apollo/mutation';
import { Link, useLocation, useParams } from 'react-router-dom';

const Session = () => {
    const location = useLocation();
    const projectName = location?.state?.projectName;
    const { id } = useParams();

    const [deleteSession, { data: deleteData, loading: deleteLoading }] =
        useMutation(DELETE_SESSION);

    const { data, loading } = useQuery(GET_SESSION, {
        variables: { id }
    });

    const onDelete = useCallback(
        (e) => {
            e.preventDefault();
            deleteSession({
                variables: { sessionId: id }
            }).then(() => {
                window.location.href = projectName
                    ? `/project/${projectName}`
                    : '/projects';
            });
        },
        [deleteSession]
    );

    const session = data?.session;

    if (loading) {
        return <Loading />;
    }

    return (
        <div className="max-w-2xl px-4 mx-auto mt-8">
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
            {Specs(session?.backlog)}
            <br />
            <br />
            <div>Machines:</div>
            <div>{ByMachine(session)}</div>
            <br />
            <br />
            <br />
            <DeleteButton
                title="Delete session"
                onClick={onDelete}
                data={deleteData}
                loading={deleteLoading}
            />
        </div>
    );
};

const Specs = (backlog) => {
    return (
        backlog?.length &&
        [...backlog]
            .sort((a, b) => b.estimatedDuration - a.estimatedDuration)
            .map((spec) => Spec(spec))
    );
};

const Spec = (spec) => {
    return (
        <div key={spec.file}>
            "{spec.file}", duration:
            {secondsToDuration(spec.estimatedDuration)}, executor:
            {spec.assignedTo}
        </div>
    );
};

const ByMachine = (session) => {
    const sessionMachines = Array.from(
        new Set(session?.backlog.map((item) => item?.assignedTo))
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

    return statsPerMachine
        .sort((a, b) => a.machine.localeCompare(b.machine))
        .map((stat) => (
            <div key={stat.machine}>
                "{stat.machine}" done in {secondsToDuration(stat.duration)}
            </div>
        ));
};

export default memo(Session);
