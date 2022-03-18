import { memo } from 'react';
import { useQuery } from '@apollo/client';
import { Link, useParams } from 'react-router-dom';
import { secondsToDuration, displayTimestamp } from '../format/displayDate';
import { defineSpecStatusTextAndColor } from '../format/specStatus';
import Loading from '../components/Loading';
import Alert from '../components/Alert';
import { GET_PROJECT } from '../apollo/query';

const Spec = () => {
    const { name, file } = useParams();

    const specFile = decodeURIComponent(file);

    const {
        data: projectData,
        loading,
        error
    } = useQuery(GET_PROJECT, {
        variables: {
            name
        }
    });

    if (error) {
        return Alert(error);
    }

    if (loading) {
        return <Loading />;
    }

    const specFileHistory = projectData?.project?.sessions
        .filter((session) => session.start > 0 && session.end > 0)
        .reduce((specFiles, session) => {
            const specFileItem = session?.backlog.find(
                (item) => item.file === specFile
            );

            specFileItem &&
                specFiles.push({
                    sessionId: session.id,
                    sessionStart: session.start,
                    sessionEnd: session.end,
                    ...specFileItem
                });
            return specFiles;
        }, [])
        .sort((a, b) => b?.sessionEnd - a?.sessionEnd);

    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            <div className="text-2xl">File "{specFile}"</div>
            <div>{specFileHistory.length} sessions</div>
            {specFileHistory.length > 0 && (
                <table className="table-auto border-collapse border border-blue-400 w-full">
                    <thead className="space-x-1">
                        <tr className="bg-blue-600 px-auto py-auto">
                            <th className="w-1/3 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    SessionID
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
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    Duration
                                </span>
                            </th>
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    Status
                                </span>
                            </th>
                            <th className="w-1/6 border border-blue-400">
                                <span className="text-gray-100 font-semibold">
                                    Duration/Session, %
                                </span>
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-gray-200">
                        {specFileHistory.map((stat) => {
                            const [statusText, bgColor] =
                                defineSpecStatusTextAndColor(stat);
                            return (
                                <tr key={stat.sessionId} className="bg-white">
                                    <td className="font-semibold border border-blue-400">
                                        <Link
                                            to={`/session/${name}/${stat.sessionId}`}
                                        >
                                            {stat.sessionId}
                                        </Link>
                                    </td>
                                    <td className="border border-blue-400">
                                        {displayTimestamp(stat.start)}
                                    </td>
                                    <td className="border border-blue-400">
                                        {displayTimestamp(stat.end)}
                                    </td>
                                    <td className="border border-blue-400">
                                        {secondsToDuration(
                                            stat.estimatedDuration
                                        )}
                                    </td>
                                    <td
                                        className={`border border-blue-400 bg-${bgColor}`}
                                    >
                                        {statusText}
                                    </td>
                                    <td className="border border-blue-400">
                                        {(
                                            (stat.estimatedDuration /
                                                (stat.sessionEnd -
                                                    stat.sessionStart)) *
                                            100
                                        ).toFixed(2)}{' '}
                                        %
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default memo(Spec);
