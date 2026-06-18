-- Revert to internal IPs
UPDATE external_services SET url='https://10.184.20.11' WHERE url='https://35.219.27.243';
UPDATE external_services SET url='https://10.184.20.13' WHERE url='https://35.219.125.94';
UPDATE external_services SET url='https://10.88.12.10' WHERE url='https://34.101.125.72';
UPDATE external_services SET url='https://10.184.20.14' WHERE url='https://34.128.98.215';
