import { memo, useCallback, useState, useEffect } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { Link, useParams } from 'react-router-dom';
import ReactPaginate from 'react-paginate';
import Loading from '../components/Loading';
import Alert from '../components/Alert';
import { displayTimestamp, secondsToDuration } from '../format/displayDate';
import { pluralize } from '../format/text';

import { DeleteButton } from '../components/DeleteButton';
import { GET_PROJECT } from '../apollo/query';
import { DELETE_PROJECT } from '../apollo/mutation';

const itemsPerPage = 15;

const Project = () => {
    const { name } = useParams();

    const [currentPage, setCurrentPage] = useState(0);
    const [pageCount, setPageCount] = useState(0);
    const [itemOffset, setItemOffset] = useState(0);
    const [itemCount, setItemCount] = useState(0);

    const pagination = {
        limit: itemsPerPage,
        offset: itemOffset
    };

    const { data, loading, error } = useQuery(GET_PROJECT, {
        variables: { name, pagination },
        fetchPolicy: 'network-only'
    });

    const handlePageClick = (event) => {
        setCurrentPage(event?.selected);
        const newOffset = (event.selected * itemsPerPage) % itemCount;
        setItemOffset(newOffset);
    };

    const [deleteProject, { data: deleteData, loading: deleteLoading }] =
        useMutation(DELETE_PROJECT);

    useEffect(() => {
        setItemCount(data?.project?.totalSessions);
        setPageCount(Math.ceil(itemCount / itemsPerPage));
    }, [itemCount, pageCount]);

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

    return error ? (
        Alert(error)
    ) : (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            <div className="text-2xl">{project?.projectName}</div>
            <div>
                {project?.sessions &&
                    Sessions(project?.projectName, project.sessions, {
                        itemCount,
                        pageCount,
                        handlePageClick,
                        currentPage
                    })}
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

const Sessions = (project, sessions, options) => {
    const orderedSessions = [...sessions].sort((a, b) => b?.start - a?.start);
    return (
        <div>
            <p>
                {options.itemCount}
                {pluralize(' session', options.itemCount)}
            </p>

            {orderedSessions?.length > 0 && (
                <div>
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
                    <div
                        id="container"
                        className="flex flex-row justify-center"
                    >
                        <ReactPaginate
                            nextLabel="next >"
                            onPageChange={options.handlePageClick}
                            pageRangeDisplayed={3}
                            marginPagesDisplayed={2}
                            forcePage={options.currentPage}
                            pageCount={options.pageCount}
                            previousLabel="< previous"
                            pageClassName="page-item"
                            pageLinkClassName="page-link"
                            previousClassName="page-item"
                            previousLinkClassName="page-link"
                            nextClassName="page-item"
                            nextLinkClassName="page-link"
                            breakLabel="..."
                            breakClassName="page-item"
                            breakLinkClassName="page-link"
                            containerClassName="pagination"
                            activeClassName="page-active"
                            disabledClassName="page-disabled"
                            renderOnZeroPageCount={null}
                        />
                    </div>
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
