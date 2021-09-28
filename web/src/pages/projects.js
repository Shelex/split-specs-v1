import { useQuery } from '@apollo/client';
import { memo } from 'react';
import { Link } from 'react-router-dom';
import { GET_PROJECTS } from '../apollo/query';
import Loading from '../components/Loading';
import Alert from '../components/Alert';

const Projects = () => {
    const { data, loading, error } = useQuery(GET_PROJECTS, {
        fetchPolicy: 'network-only'
    });

    if (loading) {
        return <Loading />;
    }

    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            {error ? (
                Alert(error)
            ) : data?.projects.length ? (
                <div>
                    <div className="text-2xl">Projects:</div>
                    <div className="grid gap-3 grid-cols-3 mt-10">
                        {data.projects.map((project) => ProjectItem(project))}
                    </div>
                </div>
            ) : (
                <ProjectsEmpty />
            )}
        </div>
    );
};

const ProjectsEmpty = () => {
    return (
        <div className="max-w-6xl px-4 mx-auto mt-8">
            No projects available. You can integrate with:
            <li>
                graphql api, schema and docs are available in
                <a
                    className="text-blue-600 mx-2"
                    href="https://split-specs.appspot.com/playground"
                >
                    graphiQL playground,
                </a>
            </li>
            <li>
                or use
                <a
                    className="text-blue-600 mx-2"
                    href="https://github.com/Shelex/split-specs-client"
                >
                    client library for js
                </a>
            </li>
            <li>
                or check
                <Link className="text-blue-600 mx-2" to="/emulate">
                    session emulation
                </Link>{' '}
                to find out how it is working :)
            </li>
        </div>
    );
};

const ProjectItem = (name) => {
    return (
        <Link to={`/project/${name}`} key={name}>
            <div className="rounded-md py-3 px-6 inline-block border-2 border-blue-600 items-center">
                <p className="align-middle break-all">{name}</p>
            </div>
        </Link>
    );
};

export default memo(Projects);
