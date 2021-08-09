import { useQuery } from '@apollo/client';
import { memo } from 'react';
import { Link } from 'react-router-dom';
import { GET_PROJECTS } from '../apollo/query';
import Loading from '../components/Loading';

const Projects = () => {
    const { data, loading } = useQuery(GET_PROJECTS);

    if (loading) {
        return <Loading />;
    }

    return (
        <div>
            {data?.projects.length ? (
                <div className="max-w-7xl px-4 mx-auto mt-8">
                    <div className="text-2xl">Available projects:</div>
                    <br />
                    <div className="grid gap-3 grid-cols-3">
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
        <p>
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
        </p>
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
