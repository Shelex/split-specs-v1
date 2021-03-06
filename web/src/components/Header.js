import { memo, useState, useCallback } from 'react';
import { Link } from 'react-router-dom';

import Logo from './Logo';
import Logout from './Logout';
import Menu from './Menu';
import Nav from './Nav';

const menuItems = [
    { name: 'Emulate Session', link: 'emulate' },
    {
        name: 'API Keys',
        link: 'apiKeys'
    }
];

const Header = ({ title }) => {
    const [isMenu, setIsMenu] = useState(false);

    const onClick = useCallback(() => {
        setIsMenu((prev) => !prev);
    }, []);

    return (
        <header className="bg-blue-800 max-w-full">
            <div className="h-14 sm:h-16 max-w-2xl mx-auto px-4 flex items-center justify-between">
                <div className="flex-1 flex items-center justify-between sm:justify-start">
                    <Menu isMenu={isMenu} onClick={onClick} />
                    <Logo title={title} />
                    <Nav items={menuItems} />
                    <Logout className="block sm:hidden" />
                </div>

                <Logout className="hidden sm:block" />
            </div>
            {isMenu && (
                <div className="sm:hidden px-4 py-2 pb-6 space-y-1">
                    {menuItems.map((item, index) => (
                        <Link
                            key={index}
                            className="menu-link"
                            to={`/${item.link}`}
                        >
                            {item.name}
                        </Link>
                    ))}
                </div>
            )}
        </header>
    );
};

export default memo(Header);
