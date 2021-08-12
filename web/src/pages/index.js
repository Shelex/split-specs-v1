import { memo } from 'react';
import { BrowserRouter as Router, Redirect, Route } from 'react-router-dom';
import { useReactiveVar } from '@apollo/client';
import { isLoggedInVar } from '../apollo';

import Layout from '../components/Layout';
import Home from './home';
import SignUp from './signup';
import Projects from './projects';
import Project from './project';
import Session from './session';
import Spec from './spec';
import Emulate from './emulate';

const Index = () => {
    const isLoggedIn = useReactiveVar(isLoggedInVar);

    return (
        <Router>
            <Layout>
                <Route exact path="/" component={Home} />
                <Route exact path="/signup" component={SignUp} />

                {isLoggedIn ? (
                    <>
                        <Route path="/projects/" component={Projects} />
                        <Route path="/project/:name" component={Project} />
                        <Route path="/spec/:name/:file" component={Spec} />
                        <Route
                            path="/session/:projectName/:id"
                            component={Session}
                        />
                        <Route path="/emulate" component={Emulate} />
                    </>
                ) : (
                    <Redirect to="/" />
                )}
            </Layout>
        </Router>
    );
};

export default memo(Index);
