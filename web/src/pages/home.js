import { useReactiveVar, useQuery } from '@apollo/client';
import { memo, useEffect } from 'react';
import { Link } from 'react-router-dom';

import { GET_PROJECTS } from '../apollo/query';

import { isLoggedInVar } from '../apollo';
import SignInForm from '../components/SignInForm';

const Home = () => {
    const isLoggedIn = useReactiveVar(isLoggedInVar);

    useEffect(() => {
        document.title = 'Split Specs';
    }, []);

    return <section>{isLoggedIn ? <ProjectsList /> : <SignInForm />}</section>;
};

const ProjectsList = () => {
    const { data } = useQuery(GET_PROJECTS, {
        fetchPolicy: 'cache-and-network'
    });

    return (
        <div>
            {data?.projects.length ? (
                <div className="max-w-2xl mt-4 p-4 mx-auto space-y-4">
                    <div className="text-2xl">Available projects:</div>
                    <ul>
                        {data.projects.map((project) => ProjectItem(project))}
                    </ul>
                </div>
            ) : (
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
            )}
        </div>
    );
};

const ProjectItem = (name) => {
    return (
        <li key={name}>
            <Link to={`project/${name}`}>{name}</Link>
        </li>
    );
};

export default memo(Home);
