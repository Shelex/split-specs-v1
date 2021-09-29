import { memo } from 'react';
import ApiKeyForm from '../components/CreateApiKeyForm';

const CreateApiKey = () => {
    return (
        <div className="h-full pt-4 sm:pt-12">
            <ApiKeyForm />
        </div>
    );
};

export default memo(CreateApiKey);
