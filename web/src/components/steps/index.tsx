import { memo, useState, useEffect } from 'react';
import type { FC } from 'react';
import { Steps } from 'antd';
import type { StepsProps } from 'antd';
import business_data from './business_data';
import { useLocation, useNavigate } from 'react-router-dom';

const { importData, importIndexData } = business_data;

const LeftSteps: FC<StepsProps> = memo(() => {
  const [currentStep, setCurrentStep] = useState(0);
  const location = useLocation();
  const navigate = useNavigate();

  const [showSteps, setShowSteps] = useState(false);
  const homePath = '/page/home';
  const [maxStep, setMaxStep] = useState(0);
  useEffect(() => {
    if (location.pathname === homePath) {
      setShowSteps(false);
    } else {
      setShowSteps(true);
      const index = importIndexData.findIndex((item) => item === location.pathname) || 0;
      if (index > maxStep) {
        setMaxStep(index);
      }
      setCurrentStep(index);
    }
  }, [location.pathname]);

  const handleChange = (current) => {
    if (current <= maxStep) {
      navigate(importIndexData[current], { replace: true });
    }
  };
  return showSteps ? (
    <Steps
      onChange={handleChange}
      className="oac-w-1/6 oac-p-4  oac-font-medium"
      direction="vertical"
      size="small"
      current={currentStep}
      items={importData}
    />
  ) : null;
});

export default LeftSteps;
