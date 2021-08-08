import { memo } from 'react';
import { useQuery } from '@apollo/client';
import { displayTimestamp, secondsToDuration } from '../dates/displayDate';
import Loading from '../components/Loading';

import { GET_SESSION } from '../apollo/query';

const Session = ({ match }) => {
    const { id } = match.params;
    const { data, loading } = useQuery(GET_SESSION, { variables: { id } });

    if (loading) {
        return <Loading />;
    }

    const session = data?.session;

    return (
        <div className="max-w-7xl px-4 mx-auto mt-8">
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
