import { memo, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useMutation, useQuery } from '@apollo/client';
import { DELETE_API_KEY } from '../apollo/mutation';
import { API_KEYS } from '../apollo/query';
import { displayTimestamp } from '../format/displayDate';
import Loading from '../components/Loading';
import Alert from '../components/Alert';

const ApiKeys = () => {
    const { data, loading, error, refetch } = useQuery(API_KEYS, {
        fetchPolicy: 'network-only'
    });

    const currentDate = new Date();

    const currentTimestamp = currentDate.valueOf() / 1000;

    let [
        deleteApiKey,
        { data: deleteData, loading: deleteLoading, error: deleteError }
    ] = useMutation(DELETE_API_KEY);

    const onDelete = useCallback(
        (e) => {
            e.preventDefault();
            deleteApiKey({
                variables: {
                    keyId: e.target.value
                }
            });
        },
        [deleteApiKey]
    );

    if (deleteLoading) {
        refetch();
    }

    return (
        <div className="h-full pt-4 sm:pt-12">
            <div className="max-w-6xl px-4 mx-auto mt-8">
                {deleteError && Alert(deleteError)}
                {error && Alert(error)}
                <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-8 my-4 rounded">
                    <Link to={`/createApiKey`}>Create api key</Link>
                </button>
                {loading || deleteLoading ? (
                    <Loading />
                ) : data?.getApiKeys.length > 0 ? (
                    <table className="table-auto border-collapse border border-blue-400 w-full">
                        <thead className="space-x-1">
                            <tr className="bg-blue-600 px-auto py-auto">
                                <th className="w-1/2 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        Name
                                    </span>
                                </th>
                                <th className="w-1/3 border border-blue-400">
                                    <span className="text-gray-100 font-semibold">
                                        ExpireAt
                                    </span>
                                </th>
                                <th className="w-1/8 border border-blue-400">
                                    <span className="text-gray-100 font-semibold"></span>
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-gray-200">
                            {data?.getApiKeys.map((apiKey) => (
                                <tr key={apiKey.id} className="bg-white">
                                    <td className="font-semibold border border-blue-400">
                                        {apiKey.name}
                                    </td>
                                    <td className="border border-blue-400 max-w-0">
                                        <p
                                            className={
                                                currentTimestamp >
                                                apiKey.expireAt
                                                    ? 'bg-red-200'
                                                    : 'bg-white'
                                            }
                                        >
                                            {displayTimestamp(apiKey.expireAt)}
                                        </p>
                                    </td>
                                    <td className="border border-blue-400">
                                        <button
                                            className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-2 rounded w-48"
                                            title="Delete project"
                                            onClick={onDelete}
                                            value={apiKey.id}
                                        >
                                            {deleteLoading ? (
                                                <Spinner />
                                            ) : (
                                                `Delete`
                                            )}
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                ) : (
                    NoApiKeys()
                )}
            </div>
        </div>
    );
};

const NoApiKeys = () => {
    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            No API Keys available.
        </div>
    );
};

export default memo(ApiKeys);
